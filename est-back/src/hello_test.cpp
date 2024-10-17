#include "hello.hpp"

#include <userver/utest/utest.hpp>

UTEST(SayHelloTo, Basic) {
  EXPECT_EQ(est_back::SayHelloTo("Developer"), "Hello, Developer!\n");
  EXPECT_EQ(est_back::SayHelloTo({}), "Hello, unknown user!\n");
}
