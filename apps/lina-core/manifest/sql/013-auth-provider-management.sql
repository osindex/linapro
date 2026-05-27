-- 013：授权管理一级菜单
-- 用途：为 OIDC / OAuth 等第三方登录 provider 插件提供统一的一级菜单
-- 入口。挂载在该菜单下的插件配置页（如 Google / Discord 登录）由各
-- 源码插件通过 plugin.yaml 中的 menus[].parent_key 引用本菜单的 menu_key
-- "auth-provider" 自动建立父子关系。
--
-- 顶级目录不依赖任何具体插件存在，即使没有安装任何 OIDC 插件，仍然在
-- 角色授权、菜单树和工作台权限治理中保留可见入口。再次执行本脚本通过
-- ON CONFLICT DO NOTHING 保持幂等。

INSERT INTO sys_menu ("parent_id", "menu_key", "name", "path", "component", "perms", "icon", "type", "sort", "visible", "status", "is_frame", "is_cache", "remark", "created_at", "updated_at")
VALUES (0, 'auth-provider', 'Authentication Providers', 'auth-provider', '', '', 'lucide:key-round', 'D', 11, 1, 1, 0, 0, '宿主稳定目录：第三方授权登录 provider 管理', NOW(), NOW())
ON CONFLICT DO NOTHING;
