cmake_minimum_required(VERSION 3.5)
project(est-back CXX)
set(CMAKE_CXX_STANDARD 20)

include(CheckIncludeFileCXX)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

file(GLOB MODELS "../api/build/est-back-cpp/model/*.cpp" "../api/build/est-back-cpp/model/*.h")
file(GLOB SERVICES "./services/*.h")

add_executable(${PROJECT_NAME} main.cpp)

find_package(absl REQUIRED)
find_package(Protobuf REQUIRED)
find_package(CURL REQUIRED)
find_package(nlohmann_json REQUIRED)

find_package(opentelemetry-cpp CONFIG REQUIRED)
target_include_directories(${PROJECT_NAME} PRIVATE ${OPENTELEMETRY_CPP_INCLUDE_DIRS})
target_link_libraries(${PROJECT_NAME} PRIVATE ${OPENTELEMETRY_CPP_LIBRARIES})

find_package(Drogon CONFIG REQUIRED)
target_link_libraries(${PROJECT_NAME} PRIVATE Drogon::Drogon)

aux_source_directory(controllers CTL_SRC)
aux_source_directory(filters FILTER_SRC)
aux_source_directory(plugins PLUGIN_SRC)


drogon_create_views(${PROJECT_NAME} ${CMAKE_CURRENT_SOURCE_DIR}/views
        ${CMAKE_CURRENT_BINARY_DIR})

target_include_directories(${PROJECT_NAME}
        PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}
        ${CMAKE_CURRENT_SOURCE_DIR}/models)

target_include_directories(${PROJECT_NAME}
        PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}
        ${CMAKE_CURRENT_SOURCE_DIR}/controllers)



target_sources(${PROJECT_NAME}
        PRIVATE
        ${SRC_DIR}
        ${CTL_SRC}
        ${FILTER_SRC}
        ${PLUGIN_SRC}
        ${SERVICES}
        ${MODELS})

add_subdirectory(test)
