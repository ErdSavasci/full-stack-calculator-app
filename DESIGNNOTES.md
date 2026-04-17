## Design Notes

1. Wanted each operation (addition, subtraction, etc.) to be separate as a microservice.
2. Backend has different modules such as config/, services/ and the root (/).
3. Backend has also shared package as /utils.
4. Frontend is built with ReactJS using Vite and TypeScript.
5. Docker compose is used to combine all these modules and start them together.
6. Frontend calculator design is inspired by below Dribble work.
7. Unit tests are written for each module.
8. Run script is written. (It builds the images one by one, and starts the app)
9. Google Gemini 3 is used to generate some files such as Dockerfile, .dockerignore, .gitignore etc.

Inspiration: https://dribbble.com/shots/17147470-Desktop-Calculator-UI-Design