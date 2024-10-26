//
// Created by Казенин Владимир on 24.10.2024.
//
#include <drogon/drogon.h>
#include "opentelemetry/exporters/ostream/span_exporter_factory.h"
#include "opentelemetry/sdk/trace/tracer_provider_factory.h"
#include "opentelemetry/sdk/trace/batch_span_processor_factory.h"
#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/trace/provider.h"
#include "opentelemetry/common/attribute_value.h"

static std::unordered_map<std::string, opentelemetry::nostd::shared_ptr<trace_api::Span>> spans;

void initApm() {
    if (getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == nullptr) {
        std::cerr << "OTEL_EXPORTER_OTLP_ENDPOINT is not set" << std::endl;
        return;
    }

    opentelemetry::sdk::trace::BatchSpanProcessorOptions bspOpts{};
    opentelemetry::exporter::otlp::OtlpHttpExporterOptions opts;
    opts.url = getenv("OTEL_EXPORTER_OTLP_ENDPOINT");

    opts.http_headers = opentelemetry::exporter::otlp::OtlpHeaders{
        {"Authorization", "Bearer " + std::string(getenv("OTEL_EXPORTER_AUTH"))},
    };

    auto exporter = opentelemetry::exporter::otlp::OtlpHttpExporterFactory::Create(opts);
    auto processor = opentelemetry::sdk::trace::BatchSpanProcessorFactory::Create(std::move(exporter), bspOpts);
    std::shared_ptr<trace_api::TracerProvider> provider =
        opentelemetry::sdk::trace::TracerProviderFactory::Create(std::move(processor));
    trace_api::Provider::SetTracerProvider(provider);

    std::cout << "otel exporter initialized" << std::endl;

    auto tracer = opentelemetry::trace::Provider::GetTracerProvider()->GetTracer("est-back");

    drogon::app().registerPreHandlingAdvice([&tracer](const drogon::HttpRequestPtr& req) {
        auto traceId = req->getHeader("trace_id");
        auto parentSpanId = req->getHeader("parent_span_id");
        std::cout << "Parent trace: " << traceId << " span: " << parentSpanId << std::endl;

        trace_api::StartSpanOptions options;

        if (traceId.size() != 32 || parentSpanId.size() != 16) {
            std::cerr << "Invalid trace_id or parent_span_id length" << std::endl;
        } else {
            // convert to opentelemetry types
            std::vector<uint8_t> traceIdVec(traceId.begin(), traceId.end());
            std::vector<u_int8_t> parentSpanIdVec(parentSpanId.begin(), parentSpanId.end());

            auto traceIdP = trace_api::TraceId(opentelemetry::nostd::span<const uint8_t, 16>(traceIdVec.data(), 16));
            auto parentSpanIdP =
                trace_api::SpanId(opentelemetry::nostd::span<const uint8_t, 8>(parentSpanIdVec.data(), 8));

            trace_api::SpanContext spanContext = trace_api::SpanContext(
                traceIdP, parentSpanIdP, trace_api::TraceFlags(trace_api::TraceFlags::kIsSampled), true);

            options.parent = spanContext;
        }

        std::map<std::string, opentelemetry::common::AttributeValue> attributes;
        attributes["service.name"] = "est-back";
        attributes["service.version"] = "0.0.1";
        attributes["deployment.environment"] = "production";
        attributes["http.method"] = req->method();
        attributes["http.url"] = req->path();

        spans[traceId] = tracer->StartSpan(req->path(), attributes, options);
    });

    drogon::app().registerPostHandlingAdvice(
        [](const drogon::HttpRequestPtr& req, const drogon::HttpResponsePtr& resp) {
            auto traceId = req->getHeader("trace_id");
            if (spans.find(traceId) != spans.end()) {
                spans[traceId]->SetAttribute("http.status_code", resp->statusCode());
                if (resp->statusCode() >= 500) {
                    spans[traceId]->SetStatus(opentelemetry::trace::StatusCode::kError);
                } else {
                    spans[traceId]->SetStatus(opentelemetry::trace::StatusCode::kOk);
                }

                std::cout << "Sending trace: " << traceId << std::endl;
                spans[traceId]->End();

                spans.erase(traceId);
            }
        });
}