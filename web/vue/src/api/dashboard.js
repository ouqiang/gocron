import httpClient from '../utils/httpClient'


export default {
    get(query, callback) {
        httpClient.get('/dashboard', {}, callback)
    }
}