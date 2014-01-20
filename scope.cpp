#include <unity/scopes/Category.h>
extern "C" {
#include "_cgo_export.h"
}
#include "scope.h"
#include "smartptr_helper.h"

using namespace unity::scopes;

ScopeAdapter::ScopeAdapter(GoInterface goscope) : goscope(goscope) {
}

int ScopeAdapter::start(std::string const &, RegistryProxy const &) {
    return VERSION;
}

void ScopeAdapter::stop() {
}

QueryBase::UPtr ScopeAdapter::create_query(std::string const &q,
                                           VariantMap const &hints) {
    QueryBase::UPtr query(new QueryAdapter(*this, q));
    return query;
}

QueryAdapter::QueryAdapter(ScopeAdapter &scope, std::string const &query)
    : scope(scope), query(query),
      cancel_channel(makeCancelChannel(), releaseCancelChannel) {
}

void QueryAdapter::cancelled() {
    sendCancelChannel(cancel_channel.get());
}

void QueryAdapter::run(ReplyProxy const &reply) {
    callScopeQuery(
        scope.goscope,
        const_cast<char*>(query.c_str()),
        const_cast<uintptr_t*>(reinterpret_cast<const uintptr_t*>(&reply)),
        cancel_channel.get());
}
