#include <drogon/drogon.h>
#include "./apm.cpp"


int main() {
    // Load config
    drogon::app().loadConfigFile("../config.json");

    // Register db client
    auto dbConfig = drogon::orm::PostgresConfig{
            .host = "postgres.est-dbs.svc.cluster.local",
            .port = 5432,
            .databaseName = "est-data",
            .username = getenv("POSTGRES_USERNAME"),
            .password = getenv("POSTGRES_PASSWORD"),
            .connectionNumber = 10,
            .name = "est-data",
            .timeout = 60,
            .autoBatch = true
    };
    drogon::app().addDbClient(dbConfig);

    // Init Apm
    initApm();

    // Start
    drogon::app().run();
    return 0;
}
