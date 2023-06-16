/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import moment from 'moment';
import bytes from 'bytes';


// BACKGROUND STATUS
export const getStatusClass = (status) => {
    if (status != 'finished' && status != 'failed' && status != 'started') {
        return 'default';
    }
    return status
}

// FORMATDATE
export const formatDate = (date) => {
    return moment(date).format('MM/DD/YYYY HH:mm')
}

//FORMATTING BYTES TO GB
export const formatByteForGB = (byte) => {
    return bytes(byte).toString('GB');
}

//ERROR API MESSAGE ERROR
export const errorApi = 'An error occurred while retrieving data from the API. Please try again.'

