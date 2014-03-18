#ifndef UNITYSCOPE_SHIM_H
#define UNITYSCOPE_SHIM_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* A typedef that can be used to represent a std::shared_ptr */

typedef uintptr_t SharedPtrData[2];

typedef void _CategorisedResult;

void run_scope(const char *scope_name, const char *runtime_config,
               void *pointer_to_iface);

/* SearchReply objects */
void init_search_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_search_reply_ptr(SharedPtrData data);

void search_reply_finished(SharedPtrData reply);
void search_reply_error(SharedPtrData reply, const char *err_string);
void search_reply_register_category(SharedPtrData reply, const char *id, const char *title, const char *icon, SharedPtrData category);
void search_reply_push(SharedPtrData reply, _CategorisedResult *result, char **error);

/* PreviewReply objects */
void init_preview_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_preview_reply_ptr(SharedPtrData data);

void preview_reply_finished(SharedPtrData reply);
void preview_reply_error(SharedPtrData reply, const char *err_string);
void preview_reply_push_widgets(SharedPtrData reply, void *gostring_array, int count, char **error);

/* Category objects */
void destroy_category_ptr(SharedPtrData data);

/* CategorisedResult objects */
_CategorisedResult *new_categorised_result(SharedPtrData category);
void destroy_categorised_result(_CategorisedResult *res);

void categorised_result_set_uri(_CategorisedResult *res, const char *uri);
void categorised_result_set_title(_CategorisedResult *res, const char *title);
void categorised_result_set_art(_CategorisedResult *res, const char *art);
void categorised_result_set_dnd_uri(_CategorisedResult *res, const char *uri);

#ifdef __cplusplus
}
#endif

#endif
