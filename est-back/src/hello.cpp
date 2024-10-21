#include "hello.hpp"
#include "../../api/build/est-back-cpp/model/BackBoardDto.h"

#include <fmt/format.h>

#include <userver/storages/postgres/cluster.hpp>
#include <userver/storages/postgres/component.hpp>
#include <userver/server/handlers/http_handler_base.hpp>
#include <userver/components/component.hpp>

namespace est_back {

namespace {

class Hello final : public userver::server::handlers::HttpHandlerBase {
public:
    static constexpr std::string_view kName = "handler-hello";

    using HttpHandlerBase::HttpHandlerBase;

    Hello(const userver::components::ComponentConfig &config,
          const userver::components::ComponentContext &component_context)
        : HttpHandlerBase(config, component_context),
          pg_cluster_(
              component_context.FindComponent<userver::components::Postgres>("postgres-db-1")
                  .GetCluster()) {
    }

    std::string HandleRequestThrow(const userver::server::http::HttpRequest &request,
                                   userver::server::request::RequestContext &) const override {
        est_back::model::BackBoardDto b;
        request.GetHttpResponse().SetContentType("application/json");
        auto res = pg_cluster_->Execute(userver::storages::postgres::ClusterHostType::kMaster,
                                        "select * from users");
        if (res.IsEmpty())
            return est_back::SayHelloTo("FUUUUUCK");
        return est_back::SayHelloTo(request.GetArg("name"));
    }

private:
    int to_string_ = 1;
    userver::storages::postgres::ClusterPtr pg_cluster_;
};

}  // namespace

std::string SayHelloTo(std::string_view name) {
    if (name.empty()) {
        name = "unknown user";
    }

    return fmt::format("Hello, {}!\n", name);
}

void AppendHello(userver::components::ComponentList &component_list) {
    component_list.Append<Hello>();
}

}  // namespace est_back
