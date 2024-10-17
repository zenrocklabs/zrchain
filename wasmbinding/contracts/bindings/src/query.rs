use cosmwasm_schema::{ 
    cw_serde, 
    QueryResponses 
};

use cosmwasm_std::{ 
    QueryRequest, 
    QuerierWrapper, 
    CustomQuery, 
    StdResult,
    Uint64,
};

#[cw_serde]
pub struct KeyringByAddressResponse {
    pub keyring: Option<Keyring>,
}

#[cw_serde]
pub struct Keyring {
    address: String,
    creator: String,
    description: String,
    admins: Box<[String]>,
    parties: Box<[String]>,
    key_req_fee: Option<Uint64>,
    sig_req_fee: Option<Uint64>,
    is_active: bool,
}

pub struct ZenRockQuerier<'a> {
    querier: &'a QuerierWrapper<'a, ZenRockQuery>,
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum ZenRockQuery {
    #[returns(KeyringByAddressResponse)]
    KeyringByAddressQuery {
        keyring_addr: String,
    },
}

impl CustomQuery for ZenRockQuery {}

impl<'a> ZenRockQuerier<'a> {
    pub fn new(querier: &'a QuerierWrapper<ZenRockQuery>) -> Self {
        ZenRockQuerier { querier }
    }

    pub fn get_keyring_by_address(
        &self,
        keyring_addr: String,
    ) -> StdResult<KeyringByAddressResponse> {
        let qry = ZenRockQuery::KeyringByAddressQuery {
            keyring_addr
        };
        let request: QueryRequest<ZenRockQuery> = ZenRockQuery::into(qry);
        self.querier.query(&request)
    }
}
