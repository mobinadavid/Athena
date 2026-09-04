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
	}, 	{
		Name:  "permission-group-delete",
		Title: datatypes.JSON(`{"fa": "حذف گروه دسترسی‌ها"}`),
	},
	{
		Name:  "blockchain-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست بلاکچین‌ها"}`),
	},
	{
		Name:  "blockchain-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات بلاکچین"}`),
	},
	{
		Name:  "blockchain-create",
		Title: datatypes.JSON(`{"fa": "ایجاد بلاکچین"}`),
	},
	{
		Name:  "blockchain-update",
		Title: datatypes.JSON(`{"fa": "ویرایش بلاکچین"}`),
	},
	{
		Name:  "blockchain-delete",
		Title: datatypes.JSON(`{"fa": "حذف بلاکچین"}`),
	},
	{
		Name:  "blockchain-explorer-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست اکسپلورر بلاکچین"}`),
	},
	{
		Name:  "blockchain-explorer-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات اکسپلورر بلاکچین"}`),
	},
	{
		Name:  "blockchain-explorer-create",
		Title: datatypes.JSON(`{"fa": "ایجاد اکسپلورر بلاکچین"}`),
	},
	{
		Name:  "blockchain-explorer-update",
		Title: datatypes.JSON(`{"fa": "ویرایش اکسپلورر بلاکچین"}`),
	},
	{
		Name:  "blockchain-explorer-delete",
		Title: datatypes.JSON(`{"fa": "حذف اکسپلورر بلاکچین"}`),
	},
	{
		Name:  "wallet-address-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-show",
		Title: datatypes.JSON(`{"fa": "نمایش اطلاعات آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-create",
		Title: datatypes.JSON(`{"fa": "ایجاد آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-update",
		Title: datatypes.JSON(`{"fa": "ویرایش آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-delete",
		Title: datatypes.JSON(`{"fa": "حذف آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-allocate",
		Title: datatypes.JSON(`{"fa": "اختصاص آدرس کیف پول"}`),
	},
	{
		Name:  "wallet-address-transactions",
		Title: datatypes.JSON(`{"fa": "نمایش تراکنش‌های آدرس کیف پول"}`),
	},
	{
		Name:  "payment-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست درخواست‌های پرداخت"}`),
	},
	{
		Name:  "payment-show",
		Title: datatypes.JSON(`{"fa": "نمایش جزئیات درخواست پرداخت"}`),
	},
	{
		Name:  "deposit-list",
		Title: datatypes.JSON(`{"fa": "نمایش لیست واریزی‌ها"}`),
	},
	{
		Name:  "deposit-show",
		Title: datatypes.JSON(`{"fa": "نمایش جزئیات واریزی"}`),
	},
	{
		Name:  "dashboard-show",
		Title: datatypes.JSON(`{"fa": "نمایش داشبورد"}`),
	},
}

func SeedPermission() {
	for _, permission := range permissions {
		database.GetInstance().GetClient().FirstOrCreate(&permission, models.PermissionModel{Name: permission.Name})
	}

	log.Println("Permission Seeder executed successfully.")
}
