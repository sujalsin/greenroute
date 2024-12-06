import React, { useEffect, useRef } from 'react';
import { GoogleMap } from '@react-google-maps/api';
import { ChargingStation, Route } from '../types/types';

interface MapProps {
    route?: Route;
    chargingStations?: ChargingStation[];
    onMapClick?: (lat: number, lng: number) => void;
}

const defaultCenter = {
    lat: 37.7749,
    lng: -122.4194,
};

const mapContainerStyle = {
    width: '100%',
    height: '500px',
};

const Map: React.FC<MapProps> = ({ route, chargingStations, onMapClick }) => {
    const mapRef = useRef<google.maps.Map>();
    const directionsRendererRef = useRef<google.maps.DirectionsRenderer>();
    const markersRef = useRef<google.maps.Marker[]>([]);

    useEffect(() => {
        // Initialize directions renderer
        if (!directionsRendererRef.current && mapRef.current) {
            directionsRendererRef.current = new google.maps.DirectionsRenderer({
                map: mapRef.current,
                suppressMarkers: true,
            });
        }

        // Clear existing markers
        markersRef.current.forEach(marker => marker.setMap(null));
        markersRef.current = [];

        if (route) {
            // Draw route
            const directionsService = new google.maps.DirectionsService();
            const origin = route.startLocation;
            const destination = route.endLocation;

            directionsService.route(
                {
                    origin: { lat: origin.latitude, lng: origin.longitude },
                    destination: { lat: destination.latitude, lng: destination.longitude },
                    travelMode: google.maps.TravelMode.DRIVING,
                },
                (result, status) => {
                    if (status === 'OK' && result) {
                        directionsRendererRef.current?.setDirections(result);
                    }
                }
            );

            // Add markers for start and end points
            markersRef.current.push(
                new google.maps.Marker({
                    position: { lat: origin.latitude, lng: origin.longitude },
                    map: mapRef.current,
                    title: 'Start',
                    icon: {
                        url: 'http://maps.google.com/mapfiles/ms/icons/green-dot.png',
                    },
                })
            );

            markersRef.current.push(
                new google.maps.Marker({
                    position: { lat: destination.latitude, lng: destination.longitude },
                    map: mapRef.current,
                    title: 'End',
                    icon: {
                        url: 'http://maps.google.com/mapfiles/ms/icons/red-dot.png',
                    },
                })
            );
        }

        // Add charging station markers
        if (chargingStations) {
            chargingStations.forEach(station => {
                markersRef.current.push(
                    new google.maps.Marker({
                        position: {
                            lat: station.addressInfo.latitude,
                            lng: station.addressInfo.longitude,
                        },
                        map: mapRef.current,
                        title: station.addressInfo.title,
                        icon: {
                            url: 'http://maps.google.com/mapfiles/ms/icons/blue-dot.png',
                        },
                    })
                );
            });
        }
    }, [route, chargingStations]);

    const handleMapLoad = (map: google.maps.Map) => {
        mapRef.current = map;
    };

    const handleClick = (e: google.maps.MapMouseEvent) => {
        if (onMapClick && e.latLng) {
            onMapClick(e.latLng.lat(), e.latLng.lng());
        }
    };

    return (
        <GoogleMap
            mapContainerStyle={mapContainerStyle}
            center={defaultCenter}
            zoom={12}
            onLoad={handleMapLoad}
            onClick={handleClick}
            options={{
                styles: [
                    {
                        featureType: 'poi',
                        elementType: 'labels',
                        stylers: [{ visibility: 'off' }],
                    },
                ],
                fullscreenControl: false,
                streetViewControl: false,
            }}
        />
    );
};

export default Map;
