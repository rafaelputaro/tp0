import logging
from common.utils import Bet, load_bets, has_won

class Lottery:
    """Allows the lottery draw to be carried out once the number of participating agencies has been reached
    """
    def __init__(self, amount_agencies: int):
        self._amount_to_trigger = amount_agencies
        self._agencies_waiting = []

    def add_agency(self, agency_id: str):
        """ Increase the number of customers waiting by one
        """
        if not self.agency_is_waiting(agency_id):
            self._agencies_waiting.append(agency_id)

    def agency_is_waiting(self, agency_id: str):
        """ Returns if agency is waiting
        """
        return self._agencies_waiting.count(agency_id) > 0

    def get_winners_from_agency(self, agency: str):
        """ Get the winners for a given agency
        """
        if len(self._agencies_waiting) < self._amount_to_trigger:
            return None
        else:
            return self._do_get_winners_from_agency(agency)

    def get_amount_agency_wating(self):
        return len(self._agencies_waiting)

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
        bets: list[Bet] = list(load_bets())
        return Lottery._apply_filter(agency, bets)