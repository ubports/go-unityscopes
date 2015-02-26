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
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void search_reply_push_filters(SharedPtrData reply, void *filters_json, void *filter_state_json, char **error) {
    try {
        Variant filters_var = Variant::deserialize_json(from_gostring(filters_json));
        Variant filter_state_var = Variant::deserialize_json(from_gostring(filter_state_json));
        Filters filters;
        for (const auto &f : filters_var.get_array()) {
            filters.emplace_back(FilterBase::deserialize(f.get_dict()));
        }
        auto filter_state = FilterState::deserialize(filter_state_var.get_dict());
        get_ptr<SearchReply>(reply)->push(filters, filter_state);
    } catch (const std::exception &e) {
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
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void preview_reply_push_attr(SharedPtrData reply, void *key, void *json_value, char **error) {
    try {
        Variant value = Variant::deserialize_json(from_gostring(json_value));
        get_ptr<PreviewReply>(reply)->push(from_gostring(key), value);
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void destroy_category_ptr(SharedPtrData data) {
    destroy_ptr<const Category>(data);
}

