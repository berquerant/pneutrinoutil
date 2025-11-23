import { DefaultApiFactory, Configuration } from "./client";
import { enableAxiosLogger } from "./log";
import axios from "axios";

const {
  API_CLIENT_TIMEOUT_MS,
  SERVER_URI,
  EXTERNAL_SERVER_URI,
} = process.env;

const axiosInstance = enableAxiosLogger(axios.create({
  baseURL: SERVER_URI,
  timeout: parseInt(API_CLIENT_TIMEOUT_MS || "3000", 10) || 3000,
}));

const configuration = new Configuration({
  basePath: SERVER_URI,
});
const defaultApi = DefaultApiFactory(configuration, undefined, axiosInstance);

const apiServerUri = EXTERNAL_SERVER_URI;

export { apiServerUri, defaultApi };
