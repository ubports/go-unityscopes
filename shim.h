#ifndef UNITYSCOPE_SHIM_H
#define UNITYSCOPE_SHIM_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* A typedef that can be used to represent a std::shared_ptr */

typedef uintptr_t SharedPtrData[2];

typedef void _CannedQuery;
typedef void _Result;
typedef void _CategorisedResult;
typedef void _SearchMetadata;
typedef void _ActionMetadata;
typedef void _ScopeBase;

void run_scope(void *scope_name, void *runtime_config,
               void *scope_config, void *pointer_to_iface,
               char **error);

/* ScopeBase objects */
char *scope_base_scope_directory(_ScopeBase *scope);
char *scope_base_cache_directory(_ScopeBase *scope);
char *scope_base_tmp_directory(_ScopeBase *scope);
char *scope_base_settings(_ScopeBase *scope);

/* SearchReply objects */
void init_search_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_search_reply_ptr(SharedPtrData data);

void search_reply_finished(SharedPtrData reply);
void search_reply_error(SharedPtrData reply, void *err_string);
void search_reply_register_category(SharedPtrData reply, void *id, void *title, void *icon, void *cat_template, SharedPtrData category);
void search_reply_register_departments(SharedPtrData reply, SharedPtrData dept);
void search_reply_push(SharedPtrData reply, _CategorisedResult *result, char **error);

/* PreviewReply objects */
void init_preview_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_preview_reply_ptr(SharedPtrData data);

void preview_reply_finished(SharedPtrData reply);
void preview_reply_error(SharedPtrData reply, void *err_string);
void preview_reply_push_widgets(SharedPtrData reply, void *gostring_array, int count, char **error);
void preview_reply_push_attr(SharedPtrData reply, void *key, void *json_value, char **error);

/* CannedQuery objects */
void destroy_canned_query(_CannedQuery *query);
_CannedQuery *new_canned_query(void *scope_id, void *query_str, void *department_id);
char *canned_query_get_scope_id(_CannedQuery *query);
char *canned_query_get_department_id(_CannedQuery *query);
char *canned_query_get_query_string(_CannedQuery *query);
void canned_query_set_department_id(_CannedQuery *query, void *department_id);
void canned_query_set_query_string(_CannedQuery *query, void *query_str);
char *canned_query_to_uri(_CannedQuery *query);

/* Category objects */
void destroy_category_ptr(SharedPtrData data);

/* CategorisedResult objects */
_Result *new_categorised_result(SharedPtrData category);
void destroy_result(_Result *res);

/* Result objects */
char *result_get_attr(_Result *res, void *attr, char **error);
void result_set_attr(_Result *res, void *attr, void *json_value, char **error);
void result_set_intercept_activation(_Result *res);

/* Department objects */
void init_department_ptr(SharedPtrData dest, SharedPtrData src);
void new_department(void *deptt_id, void *query, void *label, SharedPtrData dept, char **error);
void destroy_department_ptr(SharedPtrData data);
void department_add_subdepartment(SharedPtrData dept, SharedPtrData child);
void department_set_alternate_label(SharedPtrData dept, void *label);
void department_set_has_subdepartments(SharedPtrData dept, int subdepartments);

/* SearchMetadata objects */
void destroy_search_metadata(_SearchMetadata *metadata);
char *search_metadata_get_locale(_SearchMetadata *metadata);
char *search_metadata_get_form_factor(_SearchMetadata *metadata);
int search_metadata_get_cardinality(_SearchMetadata *metadata);
char *search_metadata_get_location(_SearchMetadata *metadata);

/* ActionMetadata objects */
void destroy_action_metadata(_ActionMetadata *metadata);
char *action_metadata_get_locale(_ActionMetadata *metadata);
char *action_metadata_get_form_factor(_ActionMetadata *metadata);

#ifdef __cplusplus
}
#endif

#endif
