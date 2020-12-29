#include <cstdarg>
#include <cstdint>
#include <cstdlib>
#include <ostream>
#include <new>

namespace ffi {

extern "C" {

const char *getCommandOutput(const char *command);

} // extern "C"

} // namespace ffi
