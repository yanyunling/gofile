rmdir /s/q .\web\dist\
cd .\web
call npm install
call npm run build
cd ..\
call go build
echo build finished
pause