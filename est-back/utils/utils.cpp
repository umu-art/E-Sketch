#include "utils.h"

namespace est_back::utils {
    bool isValidUUID(const std::string& id) {
        std::regex uuidRegex(R"(^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$)");
        return std::regex_match(id, uuidRegex);
    }

    std::string toUpper(const std::string& s) {
        std::string res;
        res.reserve(s.size());
        std::transform(s.begin(), s.end(), std::back_inserter(res), [](unsigned char c) { return std::toupper(c); });
        return res;
    }

    std::string toLower(const std::string& s) {
        std::string res;
        res.reserve(s.size());
        std::transform(s.begin(), s.end(), std::back_inserter(res), [](unsigned char c) { return std::tolower(c); });
        return res;
    }

    std::string strVectorToString(const std::vector<std::string>& v) {
        if (v.empty())
            return "";
        std::ostringstream oss;
        for (size_t i = 0; i < v.size(); ++i) {
            oss << "'" << v[i] << "'";
            if (i != v.size() - 1) {
                oss << ",";
            }
        }
        return oss.str();
    }
}  // namespace est_back::utils
