#ifndef UNITYSCOPE_HELPERS_H
#define UNITYSCOPE_HELPERS_H

#include <string>

typedef struct StrData StrData;

namespace gounityscopes {
namespace internal {

std::string from_gostring(StrData str);
std::string from_gostring(void *str);
void *as_bytes(const std::string &str, int *length);

}
}

#endif
