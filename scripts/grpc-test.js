import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../api/protos'], 'temperature.proto');

export let options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: '1m', target: 10 },
        { duration: '10s', target: 0 },
    ],
};

export default function () {
    client.connect('localhost:50051', {
        plaintext: true,
    });

    const response = client.invoke('temperature.TemperatureService/GetCurrentTemperature', {
        location: 'Rome',
    });

    check(response, {
        'status is OK': (r) => r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
}
