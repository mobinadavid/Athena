package seeders

import (
	"athena/src/database"
	"athena/src/models"
	"log"

	"gorm.io/datatypes"
)

var permissions = []models.PermissionModel{
	{
		Name:  "admin-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست ادمین‌ها"}`),
	},
	{
		Name:  "admin-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات پروفایل ادمین"}`),
	},
	{
		Name:  "admin-create",
		Title: datatypes.JSON(`{"fa": "ایجاد ادمین جدید"}`),
	},
	{
		Name:  "admin-update",
		Title: datatypes.JSON(`{"fa": "ویرایش اطلاعات ادمین"}`),
	},
	{
		Name:  "deactivate-admin",
		Title: datatypes.JSON(`{"fa": "غیر فعال/فعال کردن ادمین"}`),
	},
	{
		Name:  "admin-delete",
		Title: datatypes.JSON(`{"fa": "حذف ادمین"}`),
	},
	{
		Name:  "get-export-admin",
		Title: datatypes.JSON(`{"fa": "دریافت خروجی اکسل از لیست ادمین‌ها"}`),
	},
	{
		Name:  "user-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست کاربران"}`),
	}, {
		Name:  "get-export-user",
		Title: datatypes.JSON(`{"fa": "دریافت خروجی اکسل از لیست کاربر‌ان"}`),
	},
	{
		Name:  "user-create",
		Title: datatypes.JSON(`{"fa": "ساخت کاربر جدید"}`),
	},
	{
		Name:  "user-update",
		Title: datatypes.JSON(`{"fa": "ویرایش اطلاغات کاربر"}`),
	},
	{
		Name:  "user-get-by-nid",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات کاربر با استفاده از شناسه/کد ملی"}`),
	},
	{
		Name:  "user-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات کاربر"}`),
	}, {
		Name:  "deactivate-user",
		Title: datatypes.JSON(`{"fa": "غیر فعال/فعال کردن کاربر"}`),
	},
	{
		Name:  "role-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست نقش‌ها"}`),
	},
	{
		Name:  "role-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات نقش"}`),
	},
	{
		Name:  "role-create",
		Title: datatypes.JSON(`{"fa": "ایجاد نقش جدید"}`),
	},
	{
		Name:  "role-update",
		Title: datatypes.JSON(`{"fa": "ویرایش نقش"}`),
	},
	{
		Name:  "role-delete",
		Title: datatypes.JSON(`{"fa": "حذف نقش"}`),
	},
	{
		Name:  "permission-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست دسترسی‌ها"}`),
	}, {
		Name:  "permission-group-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست گروه دسترسی‌ها"}`),
	},
	{
		Name:  "permission-group-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات گروه دسترسی"}`),
	}, {
		Name:  "permission-group-create",
		Title: datatypes.JSON(`{"fa": "ایجاد گروه دسترسی‌ها"}`),
	}, {
		Name:  "permission-group-update",
		Title: datatypes.JSON(`{"fa": "ویرایش گروه دسترسی‌ها"}`),
	}, {
		Name:  "permission-group-delete",
		Title: datatypes.JSON(`{"fa": "حذف گروه دسترسی‌ها"}`),
	},
}

func SeedPermission() {
	for _, permission := range permissions {
		database.GetInstance().GetClient().FirstOrCreate(&permission, models.PermissionModel{Name: permission.Name})
	}

	log.Println("Permission Seeder executed successfully.")
}
