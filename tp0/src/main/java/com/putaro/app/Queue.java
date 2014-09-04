package com.putaro.app;
import com.putaro.app.SimplifiedLinkedList;

/**
 * Implementa una cola de prioridades. La misma tiene un tamaño máximo igual a MAX_QUEUE_SIZE.
 */
public class Queue 
{
	
	public final static int MAX_QUEUE_SIZE = SimplifiedLinkedList.MAX_LIST_SIZE;
	
	private SimplifiedLinkedList list = null;
	
	public Queue()
	{
		this.list = new SimplifiedLinkedList();
	}
		
	public boolean isEmpty()
	{
		return this.list.isEmpty();
	}
	
	public int size()
	{
		return this.list.getSize();
	}
	
	/**
	 * @param item Objeto que se quiere agregar en la cola.
	 * 
	 * @throws AssertionError En caso de que se trate de exceder la cantidad máxima de 
	 * elementos permitidos.
	 */
	public void add(Object item)
	{
		list.add(item);
	}
	
	/**
	 * Retorna el top de la cola.
	 * 
	 * @throws Exception En caso que la cola se encuentre vacía.
	 */
	public Object top() throws Exception
	{
		return this.list.getFirst();
	}
	
	/**
	 * Remueve el elemento top de la cola.
	 * 
	 * @throws Exception En caso que la cola se encuentre vacía.
	 */
	public void remove() throws Exception
	{
		this.list.removeFirst();
	}
}
