#include <cstring>

#include <unity/scopes/Category.h>
#include <unity/scopes/CategorisedResult.h>
#include <unity/scopes/PreviewReply.h>
#include <unity/scopes/PreviewWidget.h>
#include <unity/scopes/SearchReply.h>
#include <unity/scopes/Runtime.h>

extern "C" {
#include "_cgo_export.h"
}
#include "scope.h"
#include "smartptr_helper.h"

using namespace unity::scopes;

void run_scope(const char *scope_name, const char *runtime_config,
               void *pointer_to_iface) {
    auto runtime = Runtime::create_scope_runtime(scope_name, runtime_config);
    ScopeAdapter scope(*reinterpret_cast<GoInterface*>(pointer_to_iface));
    runtime->run_scope(&scope);
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

void search_reply_error(SharedPtrData reply, const char *err_string) {
    get_ptr<SearchReply>(reply)->error(std::make_exception_ptr(
                                     std::runtime_error(err_string)));
}

void search_reply_register_category(SharedPtrData reply, const char *id, const char *title, const char *icon, SharedPtrData category) {
    auto cat = get_ptr<SearchReply>(reply)->register_category(id, title, icon);
    init_ptr<const Category>(category, cat);
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

void preview_reply_error(SharedPtrData reply, const char *err_string) {
    get_ptr<PreviewReply>(reply)->error(std::make_exception_ptr(
                                            std::runtime_error(err_string)));
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

void categorised_result_set_uri(_CategorisedResult *res, const char *uri) {
    reinterpret_cast<CategorisedResult*>(res)->set_uri(uri);
}

void categorised_result_set_title(_CategorisedResult *res, const char *title) {
    reinterpret_cast<CategorisedResult*>(res)->set_title(title);
}

void categorised_result_set_art(_CategorisedResult *res, const char *art) {
    reinterpret_cast<CategorisedResult*>(res)->set_art(art);
}

void categorised_result_set_dnd_uri(_CategorisedResult *res, const char *uri) {
    reinterpret_cast<CategorisedResult*>(res)->set_dnd_uri(uri);
}
