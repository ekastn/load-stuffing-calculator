interface EnvConfig {
    apiBaseUrl: string;
}

const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export const config: EnvConfig = {
    apiBaseUrl,
};
