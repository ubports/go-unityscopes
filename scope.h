#ifndef UNITYSCOPE_SCOPE_H
#define UNITYSCOPE_SCOPE_H

#include <string>

#include <scopes/Reply.h>
#include <scopes/ScopeBase.h>
#include <scopes/Variant.h>

#include "_cgo_export.h"


class ScopeAdapter : public unity::api::scopes::ScopeBase
{
    friend class QueryAdapter;
public:
    ScopeAdapter(GoInterface goscope);
    virtual int start(std::string const&, unity::api::scopes::RegistryProxy const &) override;
    virtual void stop() override;
    virtual unity::api::scopes::QueryBase::UPtr create_query(std::string const &query, unity::api::scopes::VariantMap const &hints) override;

private:
    GoInterface goscope;
};

class QueryAdapter : public unity::api::scopes::QueryBase
{
public:
    QueryAdapter(ScopeAdapter &scope, std::string const &query);
    virtual void cancelled() override;
    virtual void run(unity::api::scopes::ReplyProxy const &reply) override;
private:
    const ScopeAdapter &scope;
    const std::string query;
};

#endif
