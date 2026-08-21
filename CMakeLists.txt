cmake_minimum_required(VERSION 3.20)

project(FMOS VERSION 0.1.0 LANGUAGES CXX)

add_executable(fmos
    src/main.cpp
)

target_compile_features(fmos PRIVATE cxx_std_17)

if(MSVC)
    target_compile_options(fmos PRIVATE /W4 /permissive-)
else()
    target_compile_options(fmos PRIVATE -Wall -Wextra -Wpedantic)
endif()
