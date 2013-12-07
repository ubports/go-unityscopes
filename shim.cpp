#include <scopes/Category.h>
#include <scopes/CategorisedResult.h>
#include <scopes/Reply.h>
#include <scopes/Runtime.h>

extern "C" {
#include "_cgo_export.h"
}
#include "scope.h"
#include "smartptr_helper.h"

using namespace unity::api::scopes;

void run_scope(const char *scope_name, const char *runtime_config,
               void *pointer_to_iface) {
    auto runtime = Runtime::create_scope_runtime(scope_name, runtime_config);
    ScopeAdapter scope(*reinterpret_cast<GoInterface*>(pointer_to_iface));
    runtime->run_scope(&scope);
}

void init_reply_ptr(SharedPtrData dest, SharedPtrData src) {
    std::shared_ptr<Reply> reply = get_ptr<Reply>(src);
    init_ptr<Reply>(dest, reply);
}

void destroy_reply_ptr(SharedPtrData data) {
    destroy_ptr<Reply>(data);
}

void reply_finished(SharedPtrData reply) {
    get_ptr<Reply>(reply)->finished();
}

void reply_register_category(SharedPtrData reply, const char *id, const char *title, const char *icon, SharedPtrData category) {
    auto cat = get_ptr<Reply>(reply)->register_category(id, title, icon);
    init_ptr<const Category>(category, cat);
}

void reply_push(SharedPtrData reply, _CategorisedResult *result) {
    get_ptr<Reply>(reply)->push(*reinterpret_cast<CategorisedResult*>(result));
}

void destroy_category_ptr(SharedPtrData data) {
    destroy_ptr<const Category>(data);
}

_CategorisedResult *new_categorised_result(SharedPtrData category) {
    auto cat = get_ptr<Category>(category);
    return reinterpret_cast<_CategorisedResult*>(new CategorisedResult(cat));
}

void destroy_categorised_result(_CategorisedResult *res) {
    delete reinterpret_cast<CategorisedResult*>(res);
}
