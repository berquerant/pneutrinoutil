# DefaultApi

All URIs are relative to *http://localhost:9101/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**debugGet**](#debugget) | **GET** /debug | debug info|
|[**healthGet**](#healthget) | **GET** /health | health check|
|[**procGet**](#procget) | **GET** /proc | list results|
|[**procIdConfigGet**](#procidconfigget) | **GET** /proc/{id}/config | download config|
|[**procIdDetailGet**](#prociddetailget) | **GET** /proc/{id}/detail | get process info|
|[**procIdLogGet**](#procidlogget) | **GET** /proc/{id}/log | download log|
|[**procIdMusicxmlGet**](#procidmusicxmlget) | **GET** /proc/{id}/musicxml | download musicxml|
|[**procIdWavGet**](#procidwavget) | **GET** /proc/{id}/wav | download wav|
|[**procIdWorldWavGet**](#procidworldwavget) | **GET** /proc/{id}/world_wav | download world wav|
|[**procPost**](#procpost) | **POST** /proc | start a process|
|[**procTitleGet**](#proctitleget) | **GET** /proc/title | find processes by title prefix|
|[**versionGet**](#versionget) | **GET** /version | get server version|

# **debugGet**
> HandlerSuccessResponseHandlerDebugResponseData debugGet()

debug info

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.debugGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**HandlerSuccessResponseHandlerDebugResponseData**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **healthGet**
> HandlerSuccessResponseString healthGet()

health check

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.healthGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**HandlerSuccessResponseString**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procGet**
> HandlerSuccessResponseHandlerListProcessResponseData procGet()

list results of processes

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let limit: number; //query limit; default: 5 (optional) (default to undefined)
let status: string; //process status; (pending|running|succeed|failed) (optional) (default to undefined)

const { status, data } = await apiInstance.procGet(
    limit,
    status
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **limit** | [**number**] | query limit; default: 5 | (optional) defaults to undefined|
| **status** | [**string**] | process status; (pending|running|succeed|failed) | (optional) defaults to undefined|


### Return type

**HandlerSuccessResponseHandlerListProcessResponseData**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdConfigGet**
> HandlerSuccessResponseCtlConfig procIdConfigGet()

download pneutrinoutil config as json

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdConfigGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**HandlerSuccessResponseCtlConfig**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdDetailGet**
> HandlerSuccessResponseHandlerGetDetailResponseData procIdDetailGet()

get process info

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdDetailGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**HandlerSuccessResponseHandlerGetDetailResponseData**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdLogGet**
> string procIdLogGet()

download process log file

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdLogGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdMusicxmlGet**
> string procIdMusicxmlGet()

download musicxml file

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdMusicxmlGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdWavGet**
> string procIdWavGet()

download wav file generated by pneutrinoutil

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdWavGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procIdWorldWavGet**
> string procIdWorldWavGet()

download world wav file generated by pneutrinoutil

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; //request id (default to undefined)

const { status, data } = await apiInstance.procIdWorldWavGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | request id | defaults to undefined|


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procPost**
> HandlerSuccessResponseString procPost()

start a pneutrinoutil process with given arguments

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let score: File; //musicxml (default to undefined)
let enhanceBreathiness: number; //[0, 100]%, default: 0 (optional) (default to undefined)
let formantShift: number; //default: 1.0 (optional) (default to undefined)
let inference: number; //[2, 3, 4], default: 2 (optional) (default to undefined)
let model: string; //default: MERROW (optional) (default to undefined)
let pitchShiftNsf: number; //default: 0 (optional) (default to undefined)
let pitchShiftWorld: number; //default: 0 (optional) (default to undefined)
let smoothFormant: number; //[0, 100]%, default: 0 (optional) (default to undefined)
let smoothPitch: number; //[0, 100]%, default: 0 (optional) (default to undefined)
let styleShift: number; //default: 0 (optional) (default to undefined)

const { status, data } = await apiInstance.procPost(
    score,
    enhanceBreathiness,
    formantShift,
    inference,
    model,
    pitchShiftNsf,
    pitchShiftWorld,
    smoothFormant,
    smoothPitch,
    styleShift
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **score** | [**File**] | musicxml | defaults to undefined|
| **enhanceBreathiness** | [**number**] | [0, 100]%, default: 0 | (optional) defaults to undefined|
| **formantShift** | [**number**] | default: 1.0 | (optional) defaults to undefined|
| **inference** | [**number**] | [2, 3, 4], default: 2 | (optional) defaults to undefined|
| **model** | [**string**] | default: MERROW | (optional) defaults to undefined|
| **pitchShiftNsf** | [**number**] | default: 0 | (optional) defaults to undefined|
| **pitchShiftWorld** | [**number**] | default: 0 | (optional) defaults to undefined|
| **smoothFormant** | [**number**] | [0, 100]%, default: 0 | (optional) defaults to undefined|
| **smoothPitch** | [**number**] | [0, 100]%, default: 0 | (optional) defaults to undefined|
| **styleShift** | [**number**] | default: 0 | (optional) defaults to undefined|


### Return type

**HandlerSuccessResponseString**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**202** | new process started |  * string x-request-id - request id, or just id <br>  |
|**400** | bad score |  -  |
|**413** | too big score |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **procTitleGet**
> HandlerSuccessResponseHandlerListProcessResponseData procTitleGet()

find processes by title prefix

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let prefix: string; //title prefix (default to undefined)
let limit: number; //query limit; default: 5 (optional) (default to undefined)

const { status, data } = await apiInstance.procTitleGet(
    prefix,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **prefix** | [**string**] | title prefix | defaults to undefined|
| **limit** | [**number**] | query limit; default: 5 | (optional) defaults to undefined|


### Return type

**HandlerSuccessResponseHandlerListProcessResponseData**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **versionGet**
> HandlerSuccessResponseHandlerVersionResponseData versionGet()

get server version

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.versionGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**HandlerSuccessResponseHandlerVersionResponseData**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

