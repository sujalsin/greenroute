export interface Location {
    latitude: number;
    longitude: number;
    address?: string;
}

export interface RouteSegment {
    startLocation: Location;
    endLocation: Location;
    mode: TransportMode;
    duration: number; // in seconds
    distance: number; // in meters
    co2Emission: number; // in grams
}

export interface Route {
    id?: string;
    userId?: string;
    startLocation: Location;
    endLocation: Location;
    segments: RouteSegment[];
    totalDistance: number;
    totalDuration: number;
    totalEmission: number;
    createdAt?: string;
}

export interface ChargingStation {
    id: number;
    addressInfo: {
        title: string;
        address: string;
        latitude: number;
        longitude: number;
    };
    connections: Array<{
        connectionType: {
            title: string;
        };
        powerKW: number;
    }>;
    usageType: {
        title: string;
    };
}

export interface RouteWithCharging {
    route: Route;
    chargingStations: ChargingStation[];
}

export type TransportMode = 'car' | 'bicycle' | 'public_transit' | 'walking';

export interface RoutePreferences {
    preferredModes: TransportMode[];
    avoidHighways: boolean;
    maxWalkingDistance: number; // in meters
    prioritizeEmission: boolean;
    maxTransfers: number;
}
