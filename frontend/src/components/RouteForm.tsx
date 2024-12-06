import React, { useState } from 'react';
import { RoutePreferences, TransportMode } from '../types/types';

interface RouteFormProps {
    onSubmit: (startAddress: string, endAddress: string, preferences: RoutePreferences) => void;
}

const RouteForm: React.FC<RouteFormProps> = ({ onSubmit }) => {
    const [startAddress, setStartAddress] = useState('');
    const [endAddress, setEndAddress] = useState('');
    const [preferences, setPreferences] = useState<RoutePreferences>({
        preferredModes: ['car'],
        avoidHighways: false,
        maxWalkingDistance: 1000,
        prioritizeEmission: true,
        maxTransfers: 2,
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(startAddress, endAddress, preferences);
    };

    const handleModeToggle = (mode: TransportMode) => {
        setPreferences(prev => ({
            ...prev,
            preferredModes: prev.preferredModes.includes(mode)
                ? prev.preferredModes.filter(m => m !== mode)
                : [...prev.preferredModes, mode],
        }));
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4 p-4 bg-white rounded-lg shadow">
            <div>
                <label className="block text-sm font-medium text-gray-700">Start Address</label>
                <input
                    type="text"
                    value={startAddress}
                    onChange={e => setStartAddress(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500"
                    placeholder="Enter start address"
                    required
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700">End Address</label>
                <input
                    type="text"
                    value={endAddress}
                    onChange={e => setEndAddress(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500"
                    placeholder="Enter destination address"
                    required
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700">Transport Modes</label>
                <div className="mt-2 space-x-2">
                    {(['car', 'bicycle', 'public_transit', 'walking'] as TransportMode[]).map(mode => (
                        <button
                            key={mode}
                            type="button"
                            onClick={() => handleModeToggle(mode)}
                            className={`px-3 py-1 rounded-full text-sm ${
                                preferences.preferredModes.includes(mode)
                                    ? 'bg-green-500 text-white'
                                    : 'bg-gray-200 text-gray-700'
                            }`}
                        >
                            {mode.replace('_', ' ')}
                        </button>
                    ))}
                </div>
            </div>

            <div className="space-y-2">
                <div className="flex items-center">
                    <input
                        type="checkbox"
                        checked={preferences.avoidHighways}
                        onChange={e =>
                            setPreferences(prev => ({ ...prev, avoidHighways: e.target.checked }))
                        }
                        className="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
                    />
                    <label className="ml-2 block text-sm text-gray-700">Avoid Highways</label>
                </div>

                <div className="flex items-center">
                    <input
                        type="checkbox"
                        checked={preferences.prioritizeEmission}
                        onChange={e =>
                            setPreferences(prev => ({
                                ...prev,
                                prioritizeEmission: e.target.checked,
                            }))
                        }
                        className="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
                    />
                    <label className="ml-2 block text-sm text-gray-700">
                        Prioritize Low Emissions
                    </label>
                </div>
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700">
                    Maximum Walking Distance (meters)
                </label>
                <input
                    type="number"
                    value={preferences.maxWalkingDistance}
                    onChange={e =>
                        setPreferences(prev => ({
                            ...prev,
                            maxWalkingDistance: parseInt(e.target.value),
                        }))
                    }
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500"
                    min="0"
                    step="100"
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700">
                    Maximum Transit Transfers
                </label>
                <input
                    type="number"
                    value={preferences.maxTransfers}
                    onChange={e =>
                        setPreferences(prev => ({
                            ...prev,
                            maxTransfers: parseInt(e.target.value),
                        }))
                    }
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500"
                    min="0"
                    max="5"
                />
            </div>

            <button
                type="submit"
                className="w-full bg-green-500 text-white py-2 px-4 rounded-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
            >
                Calculate Route
            </button>
        </form>
    );
};

export default RouteForm;
