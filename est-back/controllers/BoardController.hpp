#pragma once

#include <drogon/HttpController.h>
#include "../services/BoardService.h"

using namespace drogon;

namespace est_back::controller {
    class BoardController : public drogon::HttpController<BoardController> {
    public:
        METHOD_LIST_BEGIN
        ADD_METHOD_TO(BoardController::listByUserId, "/back/board/list/{userId}", Get);
        ADD_METHOD_TO(BoardController::create, "/back/board/create/{userId}", Post);
        ADD_METHOD_TO(BoardController::getByUuid, "/back/board/{boardId}", Get);
        ADD_METHOD_TO(BoardController::update, "/back/board/{boardId}", Patch);
        ADD_METHOD_TO(BoardController::deleteBoard, "/back/board/{boardId}", Delete);
        ADD_METHOD_TO(BoardController::share, "/back/board/share/{boardId}", Post);
        ADD_METHOD_TO(BoardController::updateShare, "/back/board/share/{boardId}", Patch);
        ADD_METHOD_TO(BoardController::unshare, "/back/board/share/{boardId}", Delete);
        ADD_METHOD_TO(BoardController::markAsRecent, "/back/board/recent/{userId}", Post);
        ADD_METHOD_TO(BoardController::recentsByMinute, "back/board/recents?minutes={}", Get);
        METHOD_LIST_END
    private:
        using Callback = std::function<void(const HttpResponsePtr&)>&&;
        void listByUserId(const HttpRequestPtr& req, Callback callback, std::string&& userId);
        void create(const HttpRequestPtr& req, Callback callback, std::string&& userId);
        void getByUuid(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void update(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void deleteBoard(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void share(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void updateShare(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void unshare(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void markAsRecent(const HttpRequestPtr& req, Callback callback, std::string&& userId);
        void recentsByMinute(const HttpRequestPtr& req, Callback callback, uint32_t minutes);
    };
}  // namespace est_back::controller
// namespace est- back
