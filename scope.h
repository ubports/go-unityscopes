#ifndef UNITYSCOPE_SCOPE_H
#define UNITYSCOPE_SCOPE_H

#include <memory>
#include <string>

#include <unity/scopes/SearchReply.h>
#include <unity/scopes/ScopeBase.h>
#include <unity/scopes/Variant.h>

class ScopeAdapter : public unity::scopes::ScopeBase
{
    friend class QueryAdapter;
public:
    ScopeAdapter(GoInterface goscope);
    virtual int start(std::string const&, unity::scopes::RegistryProxy const &) override;
    virtual void stop() override;
    virtual unity::scopes::SearchQueryBase::UPtr search(unity::scopes::CannedQuery const &query, unity::scopes::SearchMetadata const &hints) override;

    virtual unity::scopes::PreviewQueryBase::UPtr preview(unity::scopes::Result const& result, unity::scopes::ActionMetadata const& hints) override;

private:
    GoInterface goscope;
};

class QueryAdapter : public unity::scopes::SearchQueryBase
{
public:
    QueryAdapter(ScopeAdapter &scope, unity::scopes::CannedQuery const &query);
    virtual void cancelled() override;
    virtual void run(unity::scopes::SearchReplyProxy const &reply) override;
private:
    const ScopeAdapter &scope;
    const unity::scopes::CannedQuery query;
    std::unique_ptr<void, void(*)(GoChan)> cancel_channel;
};

#endif
