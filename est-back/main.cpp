#include <drogon/drogon.h>
#include "./apm.cpp"
#include "errors/ExceptionHandler.h"

int main() {
    // Load config
    drogon::app().loadConfigFile("../config.json");
    // Register db client
    auto dbConfig = drogon::orm::PostgresConfig{.host = "postgres.databases.svc.cluster.local",
                                                .port = 5432,
                                                .databaseName = "e-sketch",
                                                .username = getenv("POSTGRES_USERNAME"),
                                                .password = getenv("POSTGRES_PASSWORD"),
                                                .connectionNumber = 10,
                                                .name = "est-data",
                                                .timeout = 60,
                                                .autoBatch = true};
    drogon::app().addDbClient(dbConfig);

    // Init Apm
    // initApm();

    // Start
    drogon::app().setExceptionHandler(est_back::errors::customExceptionHandler);
    drogon::app().run();
    return 0;
}
