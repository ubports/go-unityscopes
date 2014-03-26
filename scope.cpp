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

SearchQueryBase::UPtr ScopeAdapter::search(CannedQuery const &q,
                                     SearchMetadata const &hints) {
    SearchQueryBase::UPtr query(new QueryAdapter(*this, q));
    return query;
}

PreviewQueryBase::UPtr ScopeAdapter::preview(Result const& result, ActionMetadata const& hints) {
    PreviewQueryBase::UPtr query(new PreviewAdapter(*this, result));
    return query;
}

QueryAdapter::QueryAdapter(ScopeAdapter &scope, CannedQuery const &query)
    : scope(scope), query(query),
      cancel_channel(makeCancelChannel(), releaseCancelChannel) {
}

void QueryAdapter::cancelled() {
    sendCancelChannel(cancel_channel.get());
}

void QueryAdapter::run(SearchReplyProxy const &reply) {
    callScopeSearch(
        scope.goscope,
        const_cast<char*>(query.query_string().c_str()),
        const_cast<uintptr_t*>(reinterpret_cast<const uintptr_t*>(&reply)),
        cancel_channel.get());
}

PreviewAdapter::PreviewAdapter(ScopeAdapter &scope, Result const &result)
    : scope(scope), result(result),
      cancel_channel(makeCancelChannel(), releaseCancelChannel) {
}

void PreviewAdapter::cancelled() {
    sendCancelChannel(cancel_channel.get());
}

void PreviewAdapter::run(PreviewReplyProxy const &reply) {
    callScopePreview(
        scope.goscope,
        reinterpret_cast<uintptr_t>(static_cast<void*>(new Result(result))),
        const_cast<uintptr_t*>(reinterpret_cast<const uintptr_t*>(&reply)),
        cancel_channel.get());
}
