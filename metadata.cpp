#include <stdexcept>
#include <cmath>
#include <cstring>

#include <unity/scopes/ActionMetadata.h>
#include <unity/scopes/SearchMetadata.h>
#include <unity/scopes/ScopeExceptions.h>

extern "C" {
#include "_cgo_export.h"
}

using namespace unity::scopes;

static void *as_bytes(const std::string &str, int *length) {
    *length = str.size();
    void *data = malloc(str.size());
    if (data == nullptr) {
        return nullptr;
    }
    memcpy(data, str.data(), str.size());
    return data;
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

void *search_metadata_get_location(_SearchMetadata *metadata, int *length) {
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
            if (!std::isfinite(value)) {
                pair.second = Variant();
            }
        }
    }
    return as_bytes(Variant(location).serialize_json(), length);
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

void *action_metadata_get_scope_data(_ActionMetadata *metadata, int *data_length) {
    const std::string data = reinterpret_cast<ActionMetadata*>(metadata)->scope_data().serialize_json();
    return as_bytes(data, data_length);
}
