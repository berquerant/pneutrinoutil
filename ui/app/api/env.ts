import { Configuration } from './client/configuration'
import { DefaultApi } from './client/api'
import { enableAxiosLogger } from './log'
import axios from 'axios'

const {
  API_CLIENT_TIMEOUT_MS,
  SERVER_URI,
} = process.env

const axiosInstance = enableAxiosLogger(axios.create({
  baseURL: SERVER_URI,
  timeout: parseInt(API_CLIENT_TIMEOUT_MS || "3000", 10) || 3000,
}))

const configuration = new Configuration({
  basePath: SERVER_URI,
})
const defaultApi = new DefaultApi(configuration, undefined, axiosInstance)

const apiServerUri = SERVER_URI

export {
  defaultApi,
  apiServerUri,
};
