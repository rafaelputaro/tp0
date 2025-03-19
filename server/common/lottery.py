import logging
from common.utils import Bet, load_bets, has_won, store_bets

KEY_AGENCIES_WAITING='agencies_waiting'

class Lottery:
    """Allows the lottery draw to be carried out once the number of participating agencies has been reached. It 
    is an intermediary for access to the betting database
    """
    def __init__(self, shared_data: dict, amount_agencies: int):
        self._amount_to_trigger = amount_agencies
        self._shared_data = shared_data

    def add_agency(self, agency_id: str):
        """ Increase the number of customers waiting by one
        """
        if not self.agency_is_waiting(agency_id):
            self._shared_data[KEY_AGENCIES_WAITING] = self._shared_data[KEY_AGENCIES_WAITING] + [agency_id]
            logging.debug(f'action: add_agency_waiting | result: success | msg: Agency ID is {agency_id} and {self.get_amount_agency_wating()} agencies are waiting.')

    def agency_is_waiting(self, agency_id: str):
        """ Returns if agency is waiting
        """
        return agency_id in self._shared_data[KEY_AGENCIES_WAITING]

    def get_winners_from_agency(self, agency: str):
        """ Get the winners for a given agency
        """
        if self.get_amount_agency_wating() < self._amount_to_trigger:
            return None
        else:
            return self._do_get_winners_from_agency(agency)

    def get_amount_agency_wating(self):
        return len(self._shared_data[KEY_AGENCIES_WAITING])

    @classmethod
    def _is_from_agency(cls, agency: int, bet: Bet):
        return bet.agency == agency

    @classmethod
    def _apply_filter(cls, agency: int, bets: list[Bet]):
        # Filters
        agency_filter = lambda bet: Lottery._is_from_agency(int(agency), bet)
        winners_filter = lambda bet: has_won(bet)
        # Apply filters
        from_this_agency = list(filter(agency_filter, bets))        
        return list(filter(winners_filter, from_this_agency))

    def _do_get_winners_from_agency(self, agency: str):
        bets: list[Bet] = list(self.load_bets())
        return Lottery._apply_filter(agency, bets)
    
    def load_bets(cls) -> list[Bet]:
        """
            Returns bets from csv
        """
        return load_bets()

    def store_bets(cls, bets: list[Bet]) -> None:
        """
            Store bets in csv
        """
        store_bets(bets)