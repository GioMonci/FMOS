#include <iostream>
#include <string_view>

namespace {

void print_usage()
{
    std::cout << "FMOS (Fix My OS)\n"
              << "\n"
              << "Usage:\n"
              << "  fmos help    Show this help message.\n";
}

} // namespace

int main(int argc, char* argv[])
{
    if (argc == 1) {
        print_usage();
        return 0;
    }

    const std::string_view command{argv[1]};
    if (command == "help" || command == "--help" || command == "-h") {
        print_usage();
        return 0;
    }

    std::cerr << "Unknown command: " << command << '\n';
    print_usage();
    return 1;
}
