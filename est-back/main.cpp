#include <drogon/drogon.h>
#include "PingController.hpp"

int main() {
    drogon::app().loadConfigFile("../config.yaml");
    drogon::app().run();
    return 0;
}
