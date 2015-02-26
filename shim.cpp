#include <cstring>

#include <unity/scopes/Category.h>
#include <unity/scopes/CategorisedResult.h>
#include <unity/scopes/PreviewReply.h>
#include <unity/scopes/PreviewWidget.h>
#include <unity/scopes/Runtime.h>
#include <unity/scopes/SearchReply.h>
#include <unity/scopes/ScopeExceptions.h>

extern "C" {
#include "_cgo_export.h"
}
#include "scope.h"
#include "smartptr_helper.h"

using namespace unity::scopes;

static std::string from_gostring(void *str) {
    GoString *s = static_cast<GoString*>(str);
    return std::string(s->p, s->n);
}

static void *as_bytes(const std::string &str, int *length) {
    *length = str.size();
    void *data = malloc(str.size());
    if (data == nullptr) {
        return nullptr;
    }
    memcpy(data, str.data(), str.size());
    return data;
}

void run_scope(void *scope_name, void *runtime_config, void *scope_config,
               void *pointer_to_iface, char **error) {
    try {
        auto runtime = Runtime::create_scope_runtime(
            from_gostring(scope_name), from_gostring(runtime_config));
        ScopeAdapter scope(*reinterpret_cast<GoInterface*>(pointer_to_iface));
        runtime->run_scope(&scope, from_gostring(scope_config));
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

char *scope_base_scope_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->scope_directory().c_str());
}

char *scope_base_cache_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->cache_directory().c_str());
}

char *scope_base_tmp_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->tmp_directory().c_str());
}

void *scope_base_settings(_ScopeBase *scope, int *length) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    Variant settings(s->settings());
    return as_bytes(settings.serialize_json(), length);
}

void destroy_category_ptr(SharedPtrData data) {
    destroy_ptr<const Category>(data);
}

