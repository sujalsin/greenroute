import { Location, RoutePreferences, RouteWithCharging } from '../types/types';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1';

export const calculateRoute = async (
    start: Location,
    end: Location,
    preferences: RoutePreferences
): Promise<RouteWithCharging> => {
    const response = await fetch(`${API_BASE_URL}/routes/calculate`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            start_location: start,
            end_location: end,
            preferences,
        }),
    });

    if (!response.ok) {
        throw new Error('Failed to calculate route');
    }

    return response.json();
};

export const getUserRoutes = async (userId: string): Promise<RouteWithCharging[]> => {
    const response = await fetch(`${API_BASE_URL}/routes/user/${userId}`);

    if (!response.ok) {
        throw new Error('Failed to fetch user routes');
    }

    return response.json();
};

export const geocodeAddress = async (address: string): Promise<Location> => {
    const response = await fetch(`${API_BASE_URL}/geocode?address=${encodeURIComponent(address)}`);

    if (!response.ok) {
        throw new Error('Failed to geocode address');
    }

    return response.json();
};
