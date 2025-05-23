export const config = {
    controller: {
        ip: process.env.CONTROLLER_IP,
        port: process.env.CONTROLLER_PORT,
    },

    API_BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL  || "localhost:8080/api",
}