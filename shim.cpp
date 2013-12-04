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
    auto runtime = Runtime::create(scope_name, runtime_config);
    ScopeAdapter scope(*reinterpret_cast<GoInterface*>(pointer_to_iface));
    //runtime->run_scope(&scope);
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
