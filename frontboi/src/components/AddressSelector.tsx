import styled from 'styled-components'
import React, { useState } from 'react'
import { FormInput, FormSelect } from './Form'

export const AddressSelectorContainer = styled.div`
    display: flex;
    flex-direction: column;
`

export function AddressSelector(props: { onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void }): React.ReactElement {
    let [addresses, setAddresses] = useState([])

    async function getAddresses(e: React.ChangeEvent<HTMLInputElement>) {
        if (e.target.value.length > 3) {
            let encodedPostcode = encodeURIComponent(e.target.value)
            let res = await fetch(`https://binboi-api.fly.dev/addresses/${encodedPostcode}`)

            let deserialisedRes = await res.json()
    
            if ("Addresses" in deserialisedRes) {
                if (Array.isArray(deserialisedRes.Addresses)) {
                    setAddresses(deserialisedRes.Addresses)
                }
            }
        }
    }

    return <>
        <AddressSelectorContainer>
            <FormInput type="text" placeholder="Postcode" onChange={getAddresses} />
            <FormSelect id="addresses" name="addresses" onChange={props.onChange}>
                {addresses.length === 0 ? <option value="">Addresses not found</option> : addresses.map((address) => {
                    return <option key={address["SiteId"]} value={address["AccountSiteUprn"]}>{address["SiteShortAddress"]}</option>
                })}
            </FormSelect>
        </AddressSelectorContainer>
    </>
}