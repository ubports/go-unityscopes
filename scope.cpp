#include <scopes/Category.h>
#include <scopes/ResultItem.h>
#include "scope.h"

using namespace unity::api::scopes;

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
    : scope(scope), query(query) {
}

void QueryAdapter::cancelled() {
    // FIXME: forward cancellation to Go (channel?)
}

void QueryAdapter::run(ReplyProxy const &reply) {
    // FIXME: spawn goroutine to handle query.
}
