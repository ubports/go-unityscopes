#include <stdexcept>
#include <cstring>

#include <unity/scopes/Department.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"
#include "smartptr_helper.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

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
    destroy_ptr<Department>(data);
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
