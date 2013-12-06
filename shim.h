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

/* Reply objects */
void init_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_reply_ptr(SharedPtrData data);

void reply_finished(SharedPtrData reply);
void reply_register_category(SharedPtrData reply, const char *id, const char *title, const char *icon, SharedPtrData category);
void reply_push(SharedPtrData reply, _CategorisedResult *result);

/* Category objects */
void destroy_category_ptr(SharedPtrData data);

/* CategorisedResult objects */
_CategorisedResult *new_categorised_result(SharedPtrData category);
void destroy_categorised_result(_CategorisedResult *res);



#ifdef __cplusplus
}
#endif

#endif
