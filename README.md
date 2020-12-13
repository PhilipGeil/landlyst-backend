# Landlyst backend

Landlyst backend består af en API server skrevet i Golang
<br>
Den er sat op til at køre på port 8080

## For at få den igang
Skal der være sat 2 miljø variabler, som indeholder login til email, som den kan sende mails fra

EMAIL_AUTH_EMAIL=test@test.dk
<br>
EMAIL_AUTH_PASS=testtest

Jeg har blot benyttet min egen ZBC mail konto til dette.

Og det kan være at man lige skal ændre stien til html siderne i email/email.go

