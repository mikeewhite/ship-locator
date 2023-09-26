import React, { useState } from 'react';
import { AutoComplete } from 'antd';
import {gql, useLazyQuery} from "@apollo/client";

const SHIP_SEARCH_QUERY = gql`
    query ($filter: String!) {
        shipSearch(searchTerm: $filter) {
            mmsi
            name
        }
    }
`;

interface AutoCompleteProps {
    onSelect: (mmsi: String | null) => void;
}

const ShipAutoComplete = ({onSelect}: AutoCompleteProps) => {
    const [options, setOptions] = useState<{ value: string; label: string }[]>([]);
    const [executeSearch, { data }] = useLazyQuery(
        SHIP_SEARCH_QUERY
    );

    const handleSearch = (value: string) => {
        let options: { value: string; label: string }[] = [];

        if (!value) {
            options = [];
        } else {
            executeSearch({
                variables: {filter: value}
            });
            options = data ? data.shipSearch.map((ship: { mmsi: string; name: string; }) => ({
                value: ship.mmsi,
                label: `${ship.name}: ${ship.mmsi}`,
            })) : [];
        }
        setOptions(options);
    }

    const handleSelect = (mmsi: String) => {
        onSelect(mmsi);
    }

    return (
        <AutoComplete
            style={{ width: 250 }}
            onSearch={handleSearch}
            placeholder="Enter Ship name or MMSI..."
            options={options}
            onSelect={handleSelect}
        />
    );
}

export default ShipAutoComplete;