# /admin/post put
curl http://127.0.0.1:8080/admin/post -X PUT \
  -d "title=titlesix" \
  -d "brief=briefsix" \
  -d "content=contentsix" \
  -d "date=20250629" \
  -d "category=20358c152d5c4f3bbe10bd3929a79d1c" \
  -d "tags=[\"53483b8f534a4ce2a566f3db9fd07ea2\",\"8148aa3523b44069a7057e3fa40dafba\"]"
# /admin/post delete
curl http://127.0.0.1:8080/admin/post/ -X PUT
# /admin/post post
curl http://127.0.0.1:8080/admin/post \
  -d "title=newtitle" \
  -d "uuid=88eb49f836fb4d9dba0156f7ec24cf8e"
