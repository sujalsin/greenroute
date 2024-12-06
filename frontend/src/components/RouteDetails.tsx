import React from 'react';
import { ChargingStation, Route } from '../types/types';

interface RouteDetailsProps {
    route: Route;
    chargingStations: ChargingStation[];
}

const RouteDetails: React.FC<RouteDetailsProps> = ({ route, chargingStations }) => {
    const formatDuration = (seconds: number): string => {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        return `${hours}h ${minutes}m`;
    };

    const formatDistance = (meters: number): string => {
        const km = meters / 1000;
        return `${km.toFixed(1)} km`;
    };

    return (
        <div className="space-y-6 p-4 bg-white rounded-lg shadow">
            <div>
                <h2 className="text-xl font-semibold text-gray-800">Route Summary</h2>
                <div className="mt-2 space-y-2">
                    <div className="flex justify-between">
                        <span className="text-gray-600">Total Distance:</span>
                        <span className="font-medium">{formatDistance(route.totalDistance)}</span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-gray-600">Total Duration:</span>
                        <span className="font-medium">
                            {formatDuration(route.totalDuration)}
                        </span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-gray-600">CO2 Emissions:</span>
                        <span className="font-medium">{route.totalEmission.toFixed(1)} g</span>
                    </div>
                </div>
            </div>

            <div>
                <h3 className="text-lg font-semibold text-gray-800">Route Segments</h3>
                <div className="mt-2 space-y-4">
                    {route.segments.map((segment, index) => (
                        <div
                            key={index}
                            className="p-3 bg-gray-50 rounded-md border border-gray-200"
                        >
                            <div className="flex justify-between items-center">
                                <span className="font-medium capitalize">
                                    {segment.mode.replace('_', ' ')}
                                </span>
                                <span className="text-sm text-gray-600">
                                    {formatDistance(segment.distance)}
                                </span>
                            </div>
                            <div className="mt-1 text-sm text-gray-600">
                                {formatDuration(segment.duration)} •{' '}
                                {segment.co2Emission.toFixed(1)} g CO2
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {chargingStations.length > 0 && (
                <div>
                    <h3 className="text-lg font-semibold text-gray-800">
                        Nearby Charging Stations
                    </h3>
                    <div className="mt-2 space-y-3">
                        {chargingStations.map(station => (
                            <div
                                key={station.id}
                                className="p-3 bg-gray-50 rounded-md border border-gray-200"
                            >
                                <div className="font-medium">{station.addressInfo.title}</div>
                                <div className="text-sm text-gray-600">
                                    {station.addressInfo.address}
                                </div>
                                <div className="mt-1 text-sm">
                                    {station.connections.map((conn, i) => (
                                        <span key={i} className="mr-2">
                                            {conn.connectionType.title} ({conn.powerKW} kW)
                                        </span>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default RouteDetails;
