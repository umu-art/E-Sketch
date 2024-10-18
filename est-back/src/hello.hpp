#pragma once

#include <string>
#include <string_view>

#include <userver/components/component_list.hpp>

namespace est_back {

std::string SayHelloTo(std::string_view name);

void AppendHello(userver::components::ComponentList &component_list);

}  // namespace est_back
