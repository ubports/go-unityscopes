#include <cmath>
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

char *scope_base_settings(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    Variant settings(s->settings());
    return strdup(settings.serialize_json().c_str());
}

void init_search_reply_ptr(SharedPtrData dest, SharedPtrData src) {
    std::shared_ptr<SearchReply> reply = get_ptr<SearchReply>(src);
    init_ptr<SearchReply>(dest, reply);
}

void destroy_search_reply_ptr(SharedPtrData data) {
    destroy_ptr<SearchReply>(data);
}

void search_reply_finished(SharedPtrData reply) {
    get_ptr<SearchReply>(reply)->finished();
}

void search_reply_error(SharedPtrData reply, void *err_string) {
    get_ptr<SearchReply>(reply)->error(std::make_exception_ptr(
        std::runtime_error(from_gostring(err_string))));
}

void search_reply_register_category(SharedPtrData reply, void *id, void *title, void *icon, void *cat_template, SharedPtrData category) {
    CategoryRenderer renderer;
    std::string renderer_template = from_gostring(cat_template);
    if (!renderer_template.empty()) {
        renderer = CategoryRenderer(renderer_template);
    }
    auto cat = get_ptr<SearchReply>(reply)->register_category(from_gostring(id), from_gostring(title), from_gostring(icon), renderer);
    init_ptr<const Category>(category, cat);
}

void search_reply_register_departments(SharedPtrData reply, SharedPtrData dept) {
    get_ptr<SearchReply>(reply)->register_departments(get_ptr<Department>(dept));
}

void search_reply_push(SharedPtrData reply, _CategorisedResult *result, char **error) {
    try {
        get_ptr<SearchReply>(reply)->push(*reinterpret_cast<CategorisedResult*>(result));
    } catch (std::exception &e) {
        *error = strdup(e.what());
    }
}

void init_preview_reply_ptr(SharedPtrData dest, SharedPtrData src) {
    std::shared_ptr<PreviewReply> reply = get_ptr<PreviewReply>(src);
    init_ptr<PreviewReply>(dest, reply);
}

void destroy_preview_reply_ptr(SharedPtrData data) {
    destroy_ptr<PreviewReply>(data);
}

void preview_reply_finished(SharedPtrData reply) {
    get_ptr<PreviewReply>(reply)->finished();
}

void preview_reply_error(SharedPtrData reply, void *err_string) {
    get_ptr<PreviewReply>(reply)->error(std::make_exception_ptr(
        std::runtime_error(from_gostring(err_string))));
}

void preview_reply_push_widgets(SharedPtrData reply, void *gostring_array, int count, char **error) {
    try {
        GoString *widget_data = static_cast<GoString*>(gostring_array);
        PreviewWidgetList widgets;
        for (int i = 0; i < count; i++) {
            widgets.push_back(PreviewWidget(std::string(
                widget_data[i].p, widget_data[i].n)));
        }
        get_ptr<PreviewReply>(reply)->push(widgets);
    } catch (std::exception &e) {
        *error = strdup(e.what());
    }
}

void preview_reply_push_attr(SharedPtrData reply, void *key, void *json_value, char **error) {
    try {
        Variant value = Variant::deserialize_json(from_gostring(json_value));
        get_ptr<PreviewReply>(reply)->push(from_gostring(key), value);
    } catch (std::exception &e) {
        *error = strdup(e.what());
    }
}

void destroy_canned_query(_CannedQuery *query) {
    delete reinterpret_cast<CannedQuery*>(query);
}

_CannedQuery *new_canned_query(void *scope_id, void *query_str, void *department_id) {
    return new CannedQuery(from_gostring(scope_id),
                           from_gostring(query_str),
                           from_gostring(department_id));
}

char *canned_query_get_scope_id(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->scope_id().c_str());
}

char *canned_query_get_department_id(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->department_id().c_str());
}

char *canned_query_get_query_string(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->query_string().c_str());
}

void canned_query_set_department_id(_CannedQuery *query, void *department_id) {
    reinterpret_cast<CannedQuery*>(query)->set_department_id(from_gostring(department_id));
}

void canned_query_set_query_string(_CannedQuery *query, void *query_str) {
    reinterpret_cast<CannedQuery*>(query)->set_query_string(from_gostring(query_str));
}

char *canned_query_to_uri(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->to_uri().c_str());
}

void destroy_category_ptr(SharedPtrData data) {
    destroy_ptr<const Category>(data);
}

_Result *new_categorised_result(SharedPtrData category) {
    auto cat = get_ptr<Category>(category);
    return reinterpret_cast<_CategorisedResult*>(static_cast<Result*>(new CategorisedResult(cat)));
}

void destroy_result(_Result *res) {
    delete reinterpret_cast<Result*>(res);
}

char *result_get_attr(_Result *res, void *attr, char **error) {
    std::string json_data;
    try {
        Variant v = reinterpret_cast<Result*>(res)->value(from_gostring(attr));
        json_data = v.serialize_json();
    } catch (std::exception &e) {
        *error = strdup(e.what());
        return nullptr;
    }
    return strdup(json_data.c_str());
}

void result_set_attr(_Result *res, void *attr, void *json_value, char **error) {
    try {
        Variant v = Variant::deserialize_json(from_gostring(json_value));
        (*reinterpret_cast<Result*>(res))[from_gostring(attr)] = v;
    } catch (std::exception &e) {
        *error = strdup(e.what());
    }
}

void result_set_intercept_activation(_Result *res) {
    reinterpret_cast<Result*>(res)->set_intercept_activation();
}

/* Department objects */
void init_department_ptr(SharedPtrData dest, SharedPtrData src) {
    std::shared_ptr<Department> dept = get_ptr<Department>(src);
    init_ptr<Department>(dest, dept);
}

void new_department(void *dept_id, void *query, void *label, SharedPtrData dept, char **error) {
    try {
        auto d = Department::create(from_gostring(dept_id),
                                    *reinterpret_cast<CannedQuery*>(query),
                                    from_gostring(label));
        init_ptr<Department>(dept, std::move(d));
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void destroy_department_ptr(SharedPtrData data) {
    destroy_ptr<PreviewReply>(data);
}

void department_add_subdepartment(SharedPtrData dept, SharedPtrData child) {
    get_ptr<Department>(dept)->add_subdepartment(get_ptr<Department>(child));
}

void department_set_alternate_label(SharedPtrData dept, void *label) {
    get_ptr<Department>(dept)->set_alternate_label(from_gostring(label));
}

void department_set_has_subdepartments(SharedPtrData dept, int subdepartments) {
    get_ptr<Department>(dept)->set_has_subdepartments(subdepartments);
}

/* SearchMetadata objects */
void destroy_search_metadata(_SearchMetadata *metadata) {
    delete reinterpret_cast<SearchMetadata*>(metadata);
}

char *search_metadata_get_locale(_SearchMetadata *metadata) {
    auto m = reinterpret_cast<SearchMetadata*>(metadata);
    try {
        return strdup(m->locale().c_str());
    } catch (const NotFoundException &) {
        return nullptr;
    }
}

char *search_metadata_get_form_factor(_SearchMetadata *metadata) {
    auto m = reinterpret_cast<SearchMetadata*>(metadata);
    try {
        return strdup(m->form_factor().c_str());
    } catch (const NotFoundException &) {
        return nullptr;
    }
}

int search_metadata_get_cardinality(_SearchMetadata *metadata) {
    return reinterpret_cast<SearchMetadata*>(metadata)->cardinality();
}

char *search_metadata_get_location(_SearchMetadata *metadata) {
    auto m = reinterpret_cast<SearchMetadata*>(metadata);
    VariantMap location;
    try {
        location = m->location().serialize();
    } catch (const NotFoundException &) {
        return nullptr;
    }
    // libjsoncpp generates invalid JSON for NaN or Inf values, so
    // filter them out here.
    for (auto &pair : location) {
        if (pair.second.which() == Variant::Double) {
            double value = pair.second.get_double();
            if (!isfinite(value)) {
                pair.second = Variant();
            }
        }
    }
    return strdup(Variant(location).serialize_json().c_str());
}

/* ActionMetadata objects */
void destroy_action_metadata(_ActionMetadata *metadata) {
    delete reinterpret_cast<ActionMetadata*>(metadata);
}

char *action_metadata_get_locale(_ActionMetadata *metadata) {
    return strdup(reinterpret_cast<ActionMetadata*>(metadata)->locale().c_str());
}

char *action_metadata_get_form_factor(_ActionMetadata *metadata) {
    return strdup(reinterpret_cast<ActionMetadata*>(metadata)->form_factor().c_str());
}
