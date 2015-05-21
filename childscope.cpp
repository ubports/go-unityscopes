#include <cstring>

#include <unity/scopes/ChildScope.h>
#include <unity/scopes/ScopeMetadata.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

_ChildScope *new_child_scope(void *id, _ScopeMetadata *metadata, int enabled, void *gostring_array, int count) {
    ScopeMetadata *api_metadata = reinterpret_cast<ScopeMetadata *>(metadata);

    GoString *keyword_data = static_cast<GoString*>(gostring_array);
    std::set<std::string> keywords;
    for (int i = 0; i < count; i++) {
        keywords.emplace(std::string(keyword_data[i].p, keyword_data[i].n));
    }

    return reinterpret_cast<_ChildScope *>(new ChildScope(from_gostring(id), *api_metadata, enabled, keywords));
}

void destroy_child_scope(_ChildScope *childscope) {
    delete reinterpret_cast<ChildScope*>(childscope);
}

char *child_scope_get_id(_ChildScope *childscope) {
    return strdup(reinterpret_cast<ChildScope*>(childscope)->id.c_str());
}
